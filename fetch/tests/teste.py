import requests

# Given URL and headers
url = "https://glue-api.vivareal.com/v3/locations?q=Rua%20Bar%C3%A3o&portal=VIVAREAL&size=6&fields=neighborhood,city,street&includeFields=address.city,address.zone,address.state,address.neighborhood,address.stateAcronym,address.street,address.locationId,address.point&"
headers = {
    "User-Agent": "Mozilla/5.0",
}

# Make the GET request
response = requests.get(url, headers=headers)

# Print the response status code and content
print("Status code:", response.status_code)
print("Response content:", response.text)
