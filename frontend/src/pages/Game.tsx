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
    const [questionText, setQuestionText] = useState<string>("");
    const [nextGameState, setNextGameState] = useState<string>("Get Question");
    const [copyText, setCopyText] = useState<string>("Click to copy share link!");
    const [copied, setCopied] = useState<boolean>(false);
    const [leader, setLeader] = useState<boolean>(false);
    const [players, setPlayers] = useState<Player[]>([]);
    const [playerCount, setPlayerCount] = useState<number>(0);
    const [waitingPlayerCount, setWaitingPlayerCount] = useState<number>(0);
    const [dropdown, setDropdown] = useState<string>("");
    const [choice, setChoice] = useState<string>("");
    const [choiceCount, setChoiceCount] = useState<number>(0);
    const [wait, setWait] = useState<number>(0);
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
                setWait(wait - 1);
            } else if (message.Command === "PLAYER JOINING") {
                setPlayers((prevPlayers) => {
                    const updatedPlayers = [
                        ...prevPlayers,
                        { name: message.Body }
                    ]
                    return sortPlayers(updatedPlayers);
                });
                console.log("ADD TO PLAYER COUNT")
                if (nextGameState === "Get Question") {
                    setPlayerCount(prevPlayerCount => prevPlayerCount + 1);
                } else {
                    setWaitingPlayerCount(prevWaitingPlayerCount => prevWaitingPlayerCount + 1);
                }
            } else if (message.Command === "PLAYER LEAVING") {
                setPlayers((prevPlayers) => {
                    const updatedPlayers = prevPlayers.filter(player => player.name !== message.Body)
                    return sortPlayers(updatedPlayers);
                });
                console.log("REMOVE FROM PLAYER COUNT")
                setPlayerCount(prevPlayerCount => prevPlayerCount - 1);
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
                setPlayerCount(otherPlayers.length);
            } else if (message.Command === "GET QUESTION") {
                setQuestionText(message.Body);
                setChoice("");
                setChoiceCount(0);
                setDropdown("");
            } else if (message.Command === "REVEAL QUESTION" || message.Command === "ODD ONE LEFT") {
                setQuestionText(message.Body);
            } else if (message.Command === "REVEAL ODD ONE OUT") {
                // chancge background color
                console.log(message.Body)
            } else if (message.Command === "NEW LEADER") {
                setLeader(true);
                setNextGameState(message.Body);
            } else if (message.Command === "NEW ROUND") {
                setNextGameState(message.Body);
            } else if (message.Command === "WAIT") {
                setWait(parseInt(message.Body));
                setQuestionText("Waiting for next round to start...")
            } else if (message.Command === "CONFIRMED CHOICES") {
                setChoiceCount(parseInt(message.Body));
                console.log("choice count: " + message.Body)
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

    const handleDropdownChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        setDropdown(event.target.value);
    }

    // confirm answer choice
    const handleChoiceButton = () => {
        if (websocketRef.current && websocketRef.current.readyState === WebSocket.OPEN && dropdown !== "") {
            websocketRef.current.send("Confirm Choice");
            setChoice(dropdown);
        }
    };

    // commands from leader
    const handleLeaderButton = () => {
        if (websocketRef.current && websocketRef.current.readyState === WebSocket.OPEN) {
            websocketRef.current.send(nextGameState);
            if (nextGameState === "Get Question") {
                setNextGameState("Reveal Question");
            } else if (nextGameState === "Reveal Question") {
                setNextGameState("Reveal Odd One Out");
            } else {
                setNextGameState("Get Question");
                setPlayerCount(prevPlayerCount => prevPlayerCount + waitingPlayerCount);
                setWaitingPlayerCount(0);
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
                <div>
                    {!choice && <select value={dropdown} onChange={handleDropdownChange} disabled={!questionText}>
                        <option key={""} value={""} disabled> Select Player </option>
                        {players.map(player => (
                            <option key={player.name}>
                                {player.name}
                            </option>
                        ))}
                    </select>}
                    {choice && <p>{choice}</p>}
                </div>
                <button type="button" onClick={handleCopyLink}>
                    {copyText}
                </button>
                <div>
                    <button type="button" onClick={handleChoiceButton} disabled={!!choice || !dropdown}> Confirm Choice </button>
                    {leader && <button type="button" onClick={handleLeaderButton} disabled={playerCount < 3 || (nextGameState === "Reveal Question" && choiceCount < playerCount)}> {nextGameState} </button>}
                </div>
            </main>
        </div>
    );
}

export default Game;


