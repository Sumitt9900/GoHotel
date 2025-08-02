# GoHotel - Full-Stack Hotel Booking Application

![GoHotel Showcase](https://placehold.co/1200x600/eef2f9/1e3a8a?text=GoHotel+Showcase)

**GoHotel** is an elegant and responsive single-page web application for booking hotel rooms. It features a modern, interactive front-end and a robust back-end powered by Go (Golang). This project demonstrates a full-stack development workflow, from creating a RESTful API to building a dynamic, user-friendly interface.

---

## Features

- **Interactive Room Selection:** Browse a gallery of high-quality room photos with a clean, modern UI.
- **Dynamic Availability:** Room cards visually update to show "Available" or "Booked" status in real-time.
- **Modal-Based Booking:** A seamless booking experience where users click a room to open a dedicated booking form.
- **Full CRUD Functionality:** Create, Read, and Delete bookings with instant feedback.
- **Real-time Updates:** The "Current Bookings" list updates automatically after any booking or cancellation.
- **Responsive Design:** A luxury experience that looks great on desktops, tablets, and mobile devices.

---

## Technology Stack

| Area          | Technology                               | Purpose                                                    |
| :------------ | :--------------------------------------- | :--------------------------------------------------------- |
| **Front-End** | **HTML5, Tailwind CSS, Vanilla JavaScript** | For a responsive, modern, and interactive user interface.  |
| **Back-End** | **Go (Golang)** | To build a fast, efficient, and scalable RESTful API server. |
| **Database** | **MongoDB** | For flexible, NoSQL data storage of booking information.     |
| **API Routing** | **Gorilla Mux** | A powerful URL router and dispatcher for Go.               |
| **CORS Handling**| **`rs/cors`** | To securely enable communication between the front-end and back-end. |

---

## Getting Started

Follow these instructions to get a local copy of the project up and running.

### Prerequisites

Ensure you have the following software installed on your machine:
- [Go (version 1.18 or later)](https://go.dev/doc/install)
- [MongoDB Community Server](https://www.mongodb.com/try/download/community)

### Installation & Setup

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/your-username/gohotel.git](https://github.com/your-username/gohotel.git)
    cd gohotel
    ```

2.  **Back-End Setup:**
    - **Install Go dependencies:** Open a terminal in the project root and run:
      ```bash
      go mod tidy
      ```
    - **Start the MongoDB Server:** Make sure your MongoDB service is running in the background. (See MongoDB's official documentation for instructions specific to your OS).
    - **Run the Go Server:** In the same terminal, start the back-end server:
      ```bash
      go run main.go
      ```
      Your API server should now be running on `http://localhost:8000`.

3.  **Front-End Setup:**
    - No special setup is needed! Simply open the `index.html` file in your web browser. You can usually do this by double-clicking the file in your file explorer.

---

## Project Structure

.
├── go.mod
├── go.sum
├── index.html      # The main front-end file
├── main.go         # The Go back-end server logic
└── script.js       # The front-end JavaScript logic


---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
