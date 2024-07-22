import React from "react";
import { useParams } from "react-router-dom";
import Header from "../layout/Header";

const Game: React.FC = () => {
    const { code } = useParams();

    return (
        <div className="game-container">
            <Header />
            <main>
                <p>
                    Room Code: {code}
                </p>
            </main>
        </div>
    );
}

export default Game;
