import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { useOutletContext } from "react-router-dom";
import config from "../config"; 

const Poll = () => {
    const { userId } = useOutletContext(); 
    const [poll, setPoll] = useState(null);
    let { id } = useParams();

    useEffect(() => {
        const headers = new Headers();
        headers.append("Content-Type", "application/json");

        const requestOptions = {
            method: "GET",
            headers: headers,
        };

        fetch(`${config.backendBaseUrl}/poll/${id}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                console.log("poll", data);
                setPoll(data.data);
            })
            .catch((err) => {
                console.log(err);
            });
    }, [id]);

    if (!poll) {
        return (
            <div>
                <h2>Poll</h2>
                <p>Loading...</p>
            </div>
        );
    }

    return (
        <div className="container mt-4">
            <div className="card shadow-sm">
                <div className="card-body">
                    <h2 className="card-title">Poll Details</h2>
                    <p><strong>ID:</strong> {poll.id}</p>
                    <p><strong>Description:</strong> {poll.description}</p>
                </div>
            </div>
            <h3 className="mt-4">Options</h3>
            <ul className="list-group">
                {poll.options && poll.options.map((option, index) => (
                    <li key={index} className="list-group-item d-flex justify-content-between align-items-center">
                        <span>{`${String.fromCharCode(97 + index)}) ${option.option_text}`}</span>
                        <div>
                            <Link to={`/vote/${option.option_id}`} className="btn btn-primary btn-sm me-2">
                                View Votes ({option.vote_count})
                            </Link>
                            <button
                                className="btn btn-success btn-sm"
                                onClick={() => handleVote(option.option_id)}
                            >
                                Vote
                            </button>
                        </div>
                    </li>
                ))}
            </ul>
        </div>
    );

    function handleVote(optionId) {
        const headers = new Headers();
        headers.append("Content-Type", "application/json");

        const payload = {
            poll_id: parseInt(id, 10), 
            option_id: parseInt(optionId, 10), 
            user_id: parseInt(userId, 10), 
        };

        console.log("payload", payload);

        const requestOptions = {
            method: "POST",
            headers: headers,
            body: JSON.stringify(payload),
        };

        fetch(`${config.backendBaseUrl}/vote`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                console.log("Vote successful", data);
                alert("Your vote has been recorded!");
            })
            .catch((err) => {
                console.error("Error voting:", err);
                alert(`An error occurred while voting: ${err.message || err}`);
            });
    }
};

export default Poll;