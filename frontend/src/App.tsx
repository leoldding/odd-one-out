import Home from "./pages/Home";
import Game from "./pages/Game";
import NotFound from "./pages/NotFound";
import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";

const App: React.FC = () => {

    return (
        <>
            <Router>
                <Routes>
                    <Route path="/" element={<Home />} />
                    <Route path="/game" element={<Game />} />
                    <Route path="*" element={<NotFound />} />
                </Routes>
            </Router>
        </>
    )
}

export default App;
