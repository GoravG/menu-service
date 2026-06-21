import json
import requests
import time

json_file = "sample_data.json"

with open(json_file, "r") as file:
    data = json.load(file)

categories = {}
menu_items = {}
portion_sizes = {}

for item in data["data"]:
    categories[item["category"]] = {
        "name": item["category"],
        "description": item["category"],
    }

    menu_items[item["menu_item_name"]] = {
        "name": item["menu_item_name"],
        "description": item["description"],
        "is_vegetarian": item["is_vegetarian"],
        "available": item["available"],
        "category": item["category"],
    }

    for price in item["prices"]:
        key = (item["menu_item_name"], price["portion_size"])
        portion_sizes[key] = {
            "menu_item_name": item["menu_item_name"],
            "portion_size": price["portion_size"],
            "price": int(price["price"]),
            "currency": price["currency"],
        }

start_time = time.time()
print("Creating Sample Data")
for category in categories.values():
    requests.post("http://localhost:8080/categories", json=category)

for menu_item in menu_items.values():
    requests.post("http://localhost:8080/menu", json=menu_item)

for portion_size in portion_sizes.values():
    requests.post("http://localhost:8080/menu-price", json=portion_size)
print("Sample Data Created")
end_time = time.time()
total_time_milliseconds = round((end_time - start_time) * 1000, 2)
print(f"Time taken: {total_time_milliseconds} milliseconds")