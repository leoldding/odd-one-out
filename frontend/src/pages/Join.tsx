import React, { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import Header from "../layout/Header";

const Join: React.FC = () => {
    const { id } = useParams();
    const navigate = useNavigate();

    const [nameValue, setNameValue] = useState<string>("");

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setNameValue(event.target.value);
    };

    const handleSubmit = () => {
        navigate("/game/" + id );
    };

    return (
        <div className="join-container">
            <Header />
            <main>
                <form onSubmit={handleSubmit}>
                    <input placeholder="Enter name" value={nameValue} onChange={handleInputChange} autoFocus />
                    <button type="submit" disabled={!nameValue}>Join Room</button>
                </form>
            </main>
        </div >
    );
}

export default Join;

