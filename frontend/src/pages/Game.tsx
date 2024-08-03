import React, { useState, useEffect } from "react";
import { useParams, useNavigate, useLocation } from "react-router-dom";
import Header from "../layout/Header";

type Player = {
    name: string
}

const Game: React.FC = () => {
    const { code } = useParams();
    const navigate = useNavigate();
    const location = useLocation();

    const [buttonText, setButtonText] = useState("Click to copy share link!");
    const [copied, setCopied] = useState(false);
    const [players, setPlayers] = useState<Player[]>([]);

    const sortPlayers = (players: Player[]) => {
        return players.sort((a, b) => a.name.localeCompare(b.name))
    };

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
        const websocket = new WebSocket("ws://localhost:8080/game")

        websocket.onopen = () => {
            websocket.send(JSON.stringify({ "name": sessionStorage.getItem("name"), "gameCode": sessionStorage.getItem("gameCode"), }))
        };

        websocket.onmessage = (event) => {
            console.log(event.data)
            const message = JSON.parse(event.data);
            console.log(message)
            if (message.Command === "PLAYER JOINING") {
                setPlayers((prevPlayers) => {
                    const updatedPlayers = [
                        ...prevPlayers,
                        { name: message.Body }
                    ]
                    return sortPlayers(updatedPlayers);
                });
            } else if (message.Command === "PLAYER LEAVING") {
                setPlayers((prevPlayers) => {
                    const updatedPlayers = prevPlayers.filter(player => player.name !== message.Body)
                    return sortPlayers(updatedPlayers);
                });
            } else if (message.Command === "OTHER PLAYERS") {
                const addPlayers = (names: string[]) => {
                    setPlayers(prevPlayers => {
                        const updatedPlayers = [
                            ...prevPlayers,
                            ...names.map(name => ({ name }))
                        ]
                        return sortPlayers(updatedPlayers);
                    });
                };
                const otherPlayers = message.Body.split(",");
                addPlayers(otherPlayers);
            }
        };

        return () => {
            if (websocket) {
                websocket.close()
            }
        };

    }, [location]);

    return (
        <div className="game-container">
            <Header />
            <main>
                <select>
                    {players.map(player => (
                        <option key={player.name}>
                            {player.name}
                        </option>
                    ))}
                </select>
                <button type="button" onClick={handleClick}>
                    {buttonText}
                </button>
            </main>
        </div>
    );
}

export default Game;
