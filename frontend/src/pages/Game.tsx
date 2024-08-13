import React, { useState, useEffect, useRef } from "react";
import { useParams, useNavigate, useLocation } from "react-router-dom";
import Header from "../layout/Header";

type Player = {
    name: string
}

const Game: React.FC = () => {
    enum gameState {
        GETQUESTION = "Get Question",
        REVEALQUESTION = "Reveal Question",
        REVEALOOO = "Reveal Odd One Out"
    }

    const { code } = useParams();
    const navigate = useNavigate();
    const location = useLocation();
    const [questionText, setQuestionText] = useState<string>("");
    const [nextGameState, setNextGameState] = useState<gameState>(gameState.GETQUESTION);
    const [copyText, setCopyText] = useState<string>("Click to copy share link!");
    const [copied, setCopied] = useState<boolean>(false);
    const [leader, setLeader] = useState<boolean>(false);
    const [players, setPlayers] = useState<Player[]>([]);
    const [playerCount, setPlayerCount] = useState<number>(0);
    const [dropdown, setDropdown] = useState<string>("");
    const dropdownRef = useRef<string>("");
    const [choice, setChoice] = useState<string>("");
    const [disableButton, setDisableButton] = useState<boolean>(true);
    const [wait, setWait] = useState<boolean>(false);
    const websocketRef = useRef<WebSocket | null>(null);

    const sortPlayers = (players: Player[]) => {
        return players.sort((a, b) => a.name.localeCompare(b.name))
    };

    const playerExists = (name: string) => {
        players.forEach((player) => {
            if (player.name == name) {
                return true;
            }
        });
        return false;
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
            if (wait && message.Command == "DONE WAITING") {
                setWait(false);
                console.log("DONE WAITING");
            }
            if (!wait) {
                if (message.Command === "PLAYERS") {
                    const addPlayers = (names: string[]) => {
                        setPlayers(() => {
                            const updatedPlayers = [
                                ...names.map(name => ({ name }))
                            ]
                            return sortPlayers(updatedPlayers);
                        });
                    };
                    const otherPlayers = message.Body.split(",");
                    addPlayers(otherPlayers);
                    setPlayerCount(otherPlayers.length);
                    if (dropdownRef.current != "" && !playerExists(dropdownRef.current)) {
                        setDropdown("");
                        dropdownRef.current = "";
                    }
                } else if (message.Command === "GET QUESTION") {
                    setQuestionText(message.Body);
                    setChoice("");
                    setDropdown("");
                    setDisableButton(true)
                } else if (message.Command === "REVEAL QUESTION" || message.Command === "ODD ONE LEFT" || message.Command === "NOT ENOUGH PLAYERS") {
                    setQuestionText(message.Body);
                } else if (message.Command === "REVEAL ODD ONE OUT") {
                    // chancge background color
                    console.log(message.Body)
                    setNextGameState(gameState.GETQUESTION);
                } else if (message.Command === "NEW LEADER") {
                    setLeader(true);
                    setNextGameState(message.Body);
                } else if (message.Command === "NEW ROUND") {
                    setNextGameState(message.Body);
                } else if (message.Command === "WAIT") {
                    setWait(true);
                    setQuestionText("Waiting for next round to start...")
                } else if (message.Command === "ALL CONFIRMED") {
                    setDisableButton(false);
                }
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
        const value = event.target.value;
        setDropdown(value);
        dropdownRef.current = value;
    }

    // confirm answer choice
    const handleChoiceButton = () => {
        if (websocketRef.current && websocketRef.current.readyState === WebSocket.OPEN && dropdown !== "") {
            websocketRef.current.send("Confirm Choice");
            setChoice(dropdownRef.current);
        }
    };

    // commands from leader
    const handleLeaderButton = () => {
        if (websocketRef.current && websocketRef.current.readyState === WebSocket.OPEN) {
            websocketRef.current.send(nextGameState);
            if (nextGameState === gameState.GETQUESTION) {
                setNextGameState(gameState.REVEALQUESTION);
            } else if (nextGameState === gameState.REVEALQUESTION) {
                setNextGameState(gameState.REVEALOOO);
            } else {
                setNextGameState(gameState.GETQUESTION);
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
                    {!choice && <select value={dropdown} onChange={handleDropdownChange} disabled={!questionText || wait}>
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
                    {leader && <button type="button" onClick={handleLeaderButton} disabled={playerCount < 3 || (nextGameState === "Reveal Question" && disableButton)}> {nextGameState} </button>}
                </div>
            </main>
        </div>
    );
}

export default Game;


