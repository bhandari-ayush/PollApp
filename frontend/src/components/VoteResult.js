import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import config from "../config"; 

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

        fetch(`${config.backendBaseUrl}/option/${id}/results`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                console.log("voteResult", data);
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
            <h2 className="mb-4">Vote Data</h2>
            <div className="vote-details mb-4">
                <p><strong>Option ID:</strong> {voteResult.option_id}</p>
                <p><strong>Vote Count:</strong> {voteResult.vote_count}</p>
            </div>
            <h3 className="mb-3">Users</h3>
            <table className="table table-striped table-bordered">
                <thead className="table-dark">
                    <tr>
                        <th style={{ width: "5%" }}>#</th>
                        <th style={{ width: "45%" }} className="text-truncate">Name</th>
                        <th style={{ width: "50%" }} className="text-truncate">Email</th>
                    </tr>
                </thead>
                <tbody>
                    {voteResult.users && voteResult.users.map((user, index) => (
                        <tr key={index}>
                            <td>{index + 1}</td>
                            <td className="text-truncate">{user.name}</td>
                            <td className="text-truncate">{user.email}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}

export default VoteResult;