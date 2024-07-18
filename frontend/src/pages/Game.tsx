import React from "react";
import { useParams } from "react-router-dom";
import Header from "../layout/Header";

const Game: React.FC = () => {
    const { id } = useParams();

    return (
        <div className="game-container">
            <Header />
            <main>
                <p>
                    Room Id: {id}
                </p>
            </main>
        </div>
    );
}

export default Game;
