export async function createRoom(name: string) {
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
    const data = await response.json();
    return data.player
}

export async function joinRoom(name: string, gameCode?: string) {
    const response = await fetch("http://localhost:8080/room/join", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ name: name, gameCode: gameCode }),
    });
    if (!response.ok) {
        throw new Error("ERROR");
    }
    const data = await response.json();
    return data.player
}
