from dynaconf import Dynaconf

settings = Dynaconf(
    settings_files=["settings.toml"],
    env_switcher="ENV",
    environments=True,
    load_dotenv=True,
)
