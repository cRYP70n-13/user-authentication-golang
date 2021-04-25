# What is this for

* This is a basic RESTful API using GoLang, MongoDB, ElasticSearch and JWT

## Endpoints

1. Authentication
      * to get your authentication key you have first to create an account and to do so u have to:

      ``` javascript
          POST localhost:8080/api/v1/users/signup

          body: {
            "first_name": "Test",
            "last_name": "User",
            "email": "testUserd@test.com",
            "password": "123456",
            "phone": "0000000000"
          }
      ```

      * And then you have to send another POST request to login and get your key

      ``` javascript
            POST localhost:8080/api/v1/users/login

            body: {
              "email": "testUserd@test.com",
              "password": "123456"
            }
      ```

      * Then as a response an Authentication Token will be provided

2. Post a question
      * To post a simple question you have to put your authentication token in the header under the name: __*token*__ \
          And then you can post a question

      ``` javascript
            POST localhost:8080/api/v1/questions

            body: {
              [
                {
                  "title": "This is a simple test",
                  "content": "A content for the test"
                },
                {
                  "title": "Anothet test",
                  "content": "Another content for the second test"
                }
              ]
            }
      ```

3. Get All the questions and some basic pagination

    ```javascript
      GET localhost:8080/api/v1/questions?query=test
    ```

4. Post an answer to a question

    ```javascript
      POST localhost:8080/api/v1/answers

      body: {
        "title": "Answer title",
        "answer": "The answer content"
      }
    ```

## The GoLang, MongoDB, ElasticSearch and JWT TODO ⛓

- [x] Implement a basic GoLang server
- [x] Implement Authentication using mongodb and JWT
- [x] Link the App with MongoDB and ElasticSearch at the same time
- [x] Implement CRUD endpoints for questions
- [x] Implement the Logic behind pagination
- [x] POST answers endpoint
- [ ] Sort the questions based on the location ElasticSearch
- [ ] Sort the questoins based on the suggest terms (elastic search) ElasticSearch
