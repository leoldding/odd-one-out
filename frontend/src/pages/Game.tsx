import React, { useState, useEffect, useRef } from "react";
import { useParams, useNavigate, useLocation } from "react-router-dom";
import Header from "../layout/Header";

type Player = {
    name: string
}

const Game: React.FC = () => {
    const { code } = useParams();
    const navigate = useNavigate();
    const location = useLocation();
    const [questionText, setQuestionText] = useState("");
    const [leaderText, setLeaderText] = useState("Get Question");
    const [copyText, setCopyText] = useState("Click to copy share link!");
    const [copied, setCopied] = useState(false);
    const [leader, setLeader] = useState(false);
    const [players, setPlayers] = useState<Player[]>([]);
    const [wait, setWait] = useState(0);
    const websocketRef = useRef<WebSocket | null>(null);

    const sortPlayers = (players: Player[]) => {
        return players.sort((a, b) => a.name.localeCompare(b.name))
    };

    // changes text of copy link button
    useEffect(() => {
        let timer: number;
        if (copied) {
            timer = window.setTimeout(() => {
                setCopyText("Click to copy share link!");
                setCopied(false);
            }, 3000);
        }
        return () => window.clearTimeout(timer);
    }, [copied]);

    // check if player has name 
    useEffect(() => {
        const name = sessionStorage.getItem("name");
        if (name === null) {
            navigate("/" + code);
        }
    }, [code, navigate]);

    // websockets
    useEffect(() => {
        const websocket = new WebSocket("ws://localhost:8080/game")
        websocketRef.current = websocket

        // register player in backend
        websocket.onopen = () => {
            websocket.send(JSON.stringify({ "name": sessionStorage.getItem("name"), "gameCode": sessionStorage.getItem("gameCode"), }))
        };

        // wait on commands from backend
        websocket.onmessage = (event) => {
            const message = JSON.parse(event.data);
            if (wait != 0) {
                setWait(wait-1);
            } else if (message.Command === "PLAYER JOINING") {
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
            } else if (message.Command === "GET QUESTION" || message.Command === "REVEAL QUESTION" || message.Command === "ODD ONE LEFT") {
                setQuestionText(message.Body);
            } else if (message.Command === "REVEAL ODD ONE OUT") {
                // chancge background color
                console.log(message.Body)
            } else if (message.Command === "NEW LEADER") {
                setLeader(true)
                setLeaderText(message.Body)
            } else if (message.Command === "NEW ROUND") {
                setLeaderText(message.Body)
            } else if (message.Command === "WAIT") {
               setWait(parseInt(message.Body)) 
               setQuestionText("Waiting for next round to start...")
            }
        };

        return () => {
            if (websocket) {
                websocket.close()
            }
        };

    }, [location]);

    // copy join link
    const handleCopyLink = () => {
        navigator.clipboard.writeText("localhost:5173/" + code);
        setCopyText("Copied!");
        setCopied(true);
    };

    // commands from leader
    const handleLeaderButton = () => {
        if (websocketRef.current && websocketRef.current.readyState === WebSocket.OPEN) {
            websocketRef.current.send(leaderText);
            if (leaderText === "Get Question") {
                setLeaderText("Reveal Question");
            } else if (leaderText === "Reveal Question") {
                setLeaderText("Reveal Odd One Out");
            } else {
                setLeaderText("Get Question")
            }
        }
    }

    return (
        <div className="game-container">
            <Header />
            <main>
                <div>
                    <h2> {questionText} </h2>
                </div>
                <select>
                    {players.map(player => (
                        <option key={player.name}>
                            {player.name}
                        </option>
                    ))}
                </select>
                <button type="button" onClick={handleCopyLink}>
                    {copyText}
                </button>
                <div>
                    <button type="button"> Confirm Choice </button>
                    {leader && <button type="button" onClick={handleLeaderButton}> {leaderText} </button>} 
                </div>
            </main>
        </div>
    );
}

export default Game;
