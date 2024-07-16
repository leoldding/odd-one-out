import React from "react";

const Home: React.FC = () => {
    return (
        <div className="home-container">
            <main>
                <h1>Odd One Out</h1>
                <input placeholder="Enter name" autoFocus />
                <button type="button">Create Room</button>
            </main>
        </div>
    );
}

export default Home;
