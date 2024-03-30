import requests

# Given URL and headers
url = "https://glue-api.zapimoveis.com.br/v2/listings?portal=ZAP&includeFields=page(link(href))&business=SALE&parentId=null&listingType=USED&addressCity=Rio de Janeiro&addressZone=Zona Oeste&addressLocationId=BR>Rio de Janeiro>NULL>Rio de Janeiro>Zona Oeste>Barra da Tijuca&addressState=Rio de Janeiro&addressNeighborhood=Barra da Tijuca&addressPointLat=-23.000371&addressPointLon=-43.365895&addressType=neighborhood&categoryPage=RESULT&developmentsSize=3&superPremiumSize=3&levels=CITY"
headers = {
    # "Host": "glue-api.zapimoveis.com.br",
    "User-Agent": "Mozilla/5.0",
    # "Accept": "application/json",
    # "Accept-Language": "pt-BR,pt;q=0.8,en-US;q=0.5,en;q=0.3",
    # "Accept-Encoding": "gzip, deflate, br",
    # "Referer": "https://www.zapimoveis.com.br/",
    "x-domain": ".zapimoveis.com.br",
    # "X-DeviceId": "5f96db06-b908-430a-b559-503a59084cb3",
    # "Authorization": "Bearer undefined",
    # "Origin": "https://www.zapimoveis.com.br",
    # "Connection": "keep-alive",
    # "Sec-Fetch-Dest": "empty",
    # "Sec-Fetch-Mode": "cors",
    # "Sec-Fetch-Site": "same-site",
    # "TE": "trailers",
}

# Make the GET request
response = requests.get(url, headers=headers)

# Print the response status code and content
print("Status code:", response.status_code)
print("Response content:", response.text)
