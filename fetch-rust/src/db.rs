use sqlx::postgres::{PgConnection, PgPool, PgPoolOptions};
use sqlx::Connection;
use std::env;

use crate::cfg::Config;

pub async fn connect_db(cfg: &Config) -> PgPool {
    let db_uri = &cfg.uri;
    let db_name = env::var("POSTGRES_DB").expect("POSTGRES_DB must be set");

    let mut con: PgConnection = Connection::connect(db_uri.as_str())
        .await
        .expect("Could not connect to database");

    con.ping().await.expect("Could not ping database");

    let pool = create_database_get_pool(&mut con, db_uri, db_name.as_str()).await;
    create_table(&pool).await;

    log::debug!("Connected to database");

    pool
}

async fn create_table(pool: &PgPool) {
    let query = include_str!("create_table.sql");
    for q in query.split(';') {
        if !q.trim().is_empty() {
            log::debug!("Running query: '{}'", q);
            sqlx::query(q)
                .execute(pool)
                .await
                .expect(format!("Could not run: '{}'", q).as_str());
        }
    }
}

async fn create_database_get_pool(con: &mut PgConnection, db_uri: &str, db_name: &str) -> PgPool {
    sqlx::query(
        format!(
            "
            DO $$BEGIN
                IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = '{nm}') THEN
                    CREATE DATABASE {nm};
                END IF;
            END$$;
            ",
            nm = db_name
        )
        .as_str(),
    )
    .execute(con)
    .await
    .expect("Could not create database");

    let url = format!("{}/{}", db_uri, db_name);
    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(url.as_str())
        .await
        .expect("Could not connect to database");

    sqlx::query("SELECT 1")
        .execute(&pool)
        .await
        .expect("Could not connect to database");

    pool
}
