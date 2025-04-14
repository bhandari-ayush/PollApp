import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";

const VoteResult = () => {
    const [voteResult, setVoteResult] = useState(null);
    let { id } = useParams();

    useEffect(() => {
        const headers = new Headers();
        headers.append("Content-Type", "application/json");

        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        fetch(`http://localhost:8080/v1/option/${id}/results`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                console.log("voteResult",data);
                setVoteResult(data);
            })
            .catch(err => {
                console.log(err);
            })
    }, [id])

    if(!voteResult){
        return(
            <div>
                <h2>Vote Data</h2>
                <p>Loading...</p>
            </div>
        )
    }

    return(
        <div className="vote-result">
            <h2>Vote Data</h2>
            <div className="vote-details">
                <p><strong>Option ID:</strong> {voteResult.option_id}</p>
                <p><strong>Vote Count:</strong> {voteResult.vote_count}</p>
            </div>
            <h3>Users</h3>
            <table className="user-table">
                <thead>
                    <tr>
                        <th>#</th>
                        <th>Name</th>
                        <th>Email</th>
                    </tr>
                </thead>
                <tbody>
                    {voteResult.users && voteResult.users.map((user, index) => (
                        <tr key={index}>
                            <td>{index + 1}</td>
                            <td>{user.name}</td>
                            <td>{user.email}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}

export default VoteResult;