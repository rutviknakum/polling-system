# Online Polling System

## Objective
The Online Polling System is a web application designed to allow users to:
- Create polls with a question and multiple answer options.
- Cast votes on existing polls.
- View real-time results of polls.

This project demonstrates core concepts like CRUD operations, real-time updates, and third-party service integration, making it both practical and impressive.

---

## Features
### 1. Create Poll
Users can create a poll by providing:
- A poll question (e.g., "What's your favorite programming language?").
- Multiple answer options (e.g., Golang, Python, Java).

### 2. Vote on Poll
Users can vote for a specific option in a poll.

### 3. View Results
Users can view poll results in real-time, showing how many votes each option has received.

### 4. Real-Time Updates
Real-time updates ensure that results are automatically updated for all users viewing the poll without needing to refresh the page.

---

## APIs

### 1. Create Poll API
- **Purpose**: Allows users to create a new poll.
- **Input (JSON format)**:
  ```json
  {
    "question": "What's your favorite programming language?",
    "options": ["Golang", "Python", "Java"]
  }
  ```
- **Output**:
  ```json
  {
    "poll_id": "12345",
    "message": "Poll created successfully!"
  }
  ```

### 2. Vote API
- **Purpose**: Enables users to vote on a poll.
- **Input (JSON format)**:
  ```json
  {
    "poll_id": "12345",
    "selected_option": "Golang"
  }
  ```
- **Output**:
  ```json
  {
    "message": "Your vote has been recorded!"
  }
  ```

### 3. View Results API
- **Purpose**: Retrieves the results of a poll.
- **Input**: Query parameter `poll_id=12345`
- **Output (JSON format)**:
  ```json
  {
    "poll_id": "12345",
    "question": "What's your favorite programming language?",
    "results": {
      "Golang": 50,
      "Python": 30,
      "Java": 20
    }
  }
  ```

---

## Technologies Used

### Backend
- **Programming Language**: Golang
- **Framework**: Gin/Gorilla Mux (for API routing)

### Database
- **Database Service**: Firebase Firestore or MongoDB Atlas
  - **Collections**:
    - `Polls`: Stores poll questions and options.
    - `Votes`: Tracks votes for each poll.

### Real-Time Updates
- **Service**: Pusher API or Firebase Realtime Database

---

## Workflow
1. **Create Poll**: The user calls the Create Poll API with a question and options. A unique poll ID is returned.
2. **Vote**: The user votes on a poll using the Vote API. The selected option is recorded.
3. **View Results**: Results are fetched using the View Results API and updated in real-time using Pusher or Firebase.


## How to Run
1. **Backend**:
   - Install Golang and required dependencies.
   - Set up the database (Firebase/MongoDB) and configure credentials.
   - Run the server: `go run main.go`

---

## Conclusion
The Online Polling System is a powerful application that combines simplicity and functionality, showcasing skills in API development, database integration, and real-time updates. This project is ideal for demonstrating practical knowledge in a real-world scenario.

