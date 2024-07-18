import Home from "./pages/Home";
import Join from "./pages/Join";
import Game from "./pages/Game";
import NotFound from "./pages/NotFound";
import React from "react";
import { BrowserRouter as Router, Route, Routes, Navigate } from "react-router-dom";

const App: React.FC = () => {

    return (
        <>
            <Router>
                <Routes>
                    <Route path="/" element={<Home />} />
                    <Route path="/join/:id/:token" element={<Join />} />
                    <Route path="/game/:id" element={<Game />} />
                    <Route path="/404" element={<NotFound />} />
                    <Route path="*" element={<Navigate to="/404" />} />
                </Routes>
            </Router>
        </>
    )
}

export default App;
