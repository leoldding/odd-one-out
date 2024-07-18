import React from "react";
import Header from "../layout/Header";

const NotFound: React.FC = () => {
    return (
        <div className="not-found-container">
            <Header />
            <h1>404 - Not Found</h1>
            <p>The page you are looking for does not exist.</p>
        </div>
    );
}

export default NotFound;
