import React from "react";

const Popup = ({ message, onClose }) => {
    return (
        <div style={{
            position: "fixed",
            top: "50%",
            left: "50%",
            transform: "translate(-50%, -50%)",
            backgroundColor: "white",
            padding: "20px",
            boxShadow: "0 4px 8px rgba(0, 0, 0, 0.2)",
            zIndex: 1050,
        }}>
            <p>{message}</p>
            <button onClick={onClose}>Close</button>
        </div>
    );
};

export default Popup;
