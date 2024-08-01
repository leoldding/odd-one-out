import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import Header from "../layout/Header";

const Game: React.FC = () => {
    const { code } = useParams();
    const navigate = useNavigate();

    const [buttonText, setButtonText] = useState("Click to copy share link!");
    const [copied, setCopied] = useState(false);

    const websocket = new WebSocket("ws://localhost:8080/game")

    const handleClick = () => {
        navigator.clipboard.writeText("localhost:5173/" + code);
        setButtonText("Copied!");
        setCopied(true);
    };

    useEffect(() => {
        let timer: number;
        if (copied) {
            timer = window.setTimeout(() => {
                setButtonText("Click to copy share link!");
                setCopied(false);
            }, 3000);
        }
        return () => window.clearTimeout(timer);
    }, [copied]);

    useEffect(() => {
        const name = sessionStorage.getItem("name");
        if (name === null) {
            navigate("/" + code);
        }
    }, []);

    useEffect(() => {
        websocket.onopen = () => {
            websocket.send(JSON.stringify({"name": sessionStorage.getItem("name"), "gameCode": sessionStorage.getItem("gameCode"),}))
        };

        websocket.onmessage = (event) => { 
            console.log(event.data)
        };
    }, [websocket]);

    return (
        <div className="game-container">
            <Header />
            <main>
                <button type="button" onClick={handleClick}>
                    {buttonText}
                </button>
            </main>
        </div>
    );
}

export default Game;
