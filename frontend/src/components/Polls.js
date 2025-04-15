import { useEffect, useState } from "react";
import { Link, useOutletContext } from "react-router-dom";

const Polls = () => {
    const { userId, setAlertMessage, setAlertClassName } = useOutletContext();
    const [polls, setPolls] = useState(null);

    const formatDate = (timestamp) => {
        const date = new Date(timestamp);
        return date.toLocaleString("en-US", {
            year: "numeric",
            month: "long",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit",
            hour12: true,
        });
    };

    const deletePoll = (pollId) => {
        if (!userId) {
            setAlertMessage("You must be logged in to delete a poll.");
            setAlertClassName("alert alert-danger");
            setTimeout(() => {
                setAlertClassName("d-none"); 
            }, 2000);
            return;
        }

        if (window.confirm("Are you sure you want to delete this poll?")) {
            const headers = new Headers();
            headers.append("Content-Type", "application/json");

            const requestOptions = {
                method: "DELETE",
                headers: headers,
            };

            fetch(`http://localhost:8080/v1/poll/${pollId}`, requestOptions)
                .then((response) => {
                    if (response.ok) {
                        setPolls(polls.filter((poll) => poll.id !== pollId));
                        
                    } else {
                        console.log("Failed to delete poll.");
                        }
      
})                .catch((err) => {
                    console.log(err);
                    
                });
        }
    };

    useEffect(() => {
        const headers = new Headers();
        headers.append("Content-Type", "application/json");

        const requestOptions = {
            method: "GET",
            headers: headers,
        };

        fetch(`http://localhost:8080/v1/all/poll/`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                setPolls(data.data);
            })
            .catch((err) => {
                console.log(err);
            });
    }, []);

    if (!polls) {
        return (
            <div>
                <h2>Polls</h2>
                <p>Loading...</p>
            </div>
        );
    }

    return (
        <div className="container mt-4">
            <h2 className="mb-4 text-center">Polls</h2>
            <div className="table-responsive">
                <table className="table table-striped table-hover table-bordered">
                    <thead className="table-dark">
                        <tr>
                            <th>ID</th>
                            <th>Description</th>
                            <th>Created At</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {polls.map((poll) => (
                            <tr key={poll.id}>
                                <td>
                                    <Link to={`/poll/${poll.id}`} className="text-decoration-none">
                                        {poll.id}
                                    </Link>
                                </td>
                                <td>{poll.description}</td>
                                <td>{formatDate(poll.created_at)}</td>
                                <td>
                                    <button
                                        className="btn btn-danger btn-sm"
                                        onClick={() => deletePoll(poll.id)}
                                    >
                                        Delete
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default Polls;