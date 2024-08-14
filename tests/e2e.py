import requests
import json

# Define the base URL
BASE_URL = "http://localhost:8001/api/v1"

API_KEY = "cJGZ8L1sDcPezjOy1zacPJZxzZxrPObm2Ggs1U0V+fE=INSECURE"  # Replace with your actual API key
headers = {"X-API-Key": API_KEY, "Content-Type": "application/json"}

# Register a new user and log in to get a JWT token
def get_jwt_token():
    # User registration (assuming it doesn't need authentication)
    register_url = f"{BASE_URL}/register"
    register_data = json.dumps({"username": "testuser", "password": "securepassword"})
    register_response = requests.post(register_url, headers=headers, data=register_data)
    print("Register User:", register_response.status_code, register_response.json())
    
    # User login
    login_url = f"{BASE_URL}/login"
    login_data = json.dumps({"username": "testuser", "password": "securepassword"})
    login_response = requests.post(login_url, headers=headers, data=login_data)
    print("Login User:", login_response.status_code, login_response.json())
    
    # Extracting JWT token from login response
    jwt_token = login_response.json().get('token')
    return jwt_token

# Create a new book
def test_create_book(jwt_token):
    url = f"{BASE_URL}/books"
    auth_headers = {**headers, "Authorization": f"Bearer {jwt_token}"}
    data = json.dumps({"author": "Jane Doe", "title": "New Book Title"})
    response = requests.post(url, headers=auth_headers, data=data)
    print("Create Book:", response.status_code, response.json())
    return response.json()

# Get all books
def test_get_books(jwt_token):
    url = f"{BASE_URL}/books"
    auth_headers = {**headers, "Authorization": f"Bearer {jwt_token}"}
    response = requests.get(url, headers=auth_headers)
    print("Get Books:", response.status_code, response.json())

# Get a specific book by ID
def test_get_book(book_id, jwt_token):
    url = f"{BASE_URL}/books/{book_id}"
    auth_headers = {**headers, "Authorization": f"Bearer {jwt_token}"}
    response = requests.get(url, headers=auth_headers)
    print("Get Book:", response.status_code, response.json())

# Update a book
def test_update_book(book_id, jwt_token):
    url = f"{BASE_URL}/books/{book_id}"
    auth_headers = {**headers, "Authorization": f"Bearer {jwt_token}"}
    data = json.dumps({"author": "John Smith", "title": "Updated Book Title"})
    response = requests.put(url, headers=auth_headers, data=data)
    print("Update Book:", response.status_code, response.json())

# Delete a book
def test_delete_book(book_id, jwt_token):
    url = f"{BASE_URL}/books/{book_id}"
    auth_headers = {**headers, "Authorization": f"Bearer {jwt_token}"}
    response = requests.delete(url, headers=auth_headers)
    print("Delete Book:", response.status_code)

# Run tests
def run_tests():
    # Authenticate and obtain JWT
    jwt_token = get_jwt_token()
    if jwt_token:
        # Create a book and use its ID for further tests
        book = test_create_book(jwt_token)
        book_id = book['data']['id']

        # Perform operations
        test_get_books(jwt_token)
        test_get_book(book_id, jwt_token)
        test_update_book(book_id, jwt_token)
        test_delete_book(book_id, jwt_token)

if __name__ == "__main__":
    run_tests()
