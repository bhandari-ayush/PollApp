import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

const Polls = () => {
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

    useEffect( () => {
        const headers = new Headers();
        headers.append("Content-Type", "application/json");

        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        fetch(`http://localhost:8080/v1/all/poll/`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                setPolls(data.data);
            })
            .catch(err => {
                console.log(err);
            })

    }, []);

    if(!polls){
        return(
            <div>
                <h2>Polls</h2>
                <p>Loading...</p>
            </div>
        )
    }

    return(
        <div className="container mt-4">
        <h2 className="mb-4 text-center">Polls</h2>
        <div className="table-responsive">
            <table className="table table-striped table-hover table-bordered">
                <thead className="table-dark">
                    <tr>
                        <th>ID</th>
                        <th>Description</th>
                        <th>Created At</th>
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
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    </div>
    )
}

export default Polls;