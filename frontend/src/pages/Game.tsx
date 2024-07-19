import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import Header from "../layout/Header";

const Game: React.FC = () => {
    const { id } = useParams();
    const [buttonText, setButtonText] = useState("Click to copy share link!");
    const [copied, setCopied] = useState(false);

    const handleClick = () => {
        navigator.clipboard.writeText("localhost:5173/join/tempID/tempToken");
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

    return (
        <div className="game-container">
            <Header />
            <main>
                <p>
                    Room Id: {id}
                </p>
                <button type="button" onClick={handleClick}>
                    {buttonText}
                </button>
            </main>
        </div>
    );
}

export default Game;
