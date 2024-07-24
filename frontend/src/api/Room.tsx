import { Player } from "../models/Player";

export async function createRoom(name: string): Promise<Player> {
    const response = await fetch("http://localhost:8080/room/create", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ name: name }),
    });
    if (!response.ok) {
        throw new Error("ERROR");
    }
    const data: Player = await response.json();
    return data
}

