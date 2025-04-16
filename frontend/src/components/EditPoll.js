import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import config from "../config"; 

const EditPoll = () => {
    const [description, setDescription] = useState("");
    const [options, setOptions] = useState([{ option_text: "" }]); // Array of options
    const [error, setError] = useState("");
    const navigate = useNavigate();
    const { id } = useParams(); 

    const handleOptionChange = (index, value) => {
        const updatedOptions = [...options];
        updatedOptions[index].option_text = value;
        setOptions(updatedOptions);
    };

    const addOption = () => {
        setOptions([...options, { option_text: "" }]);
    };

    const removeOption = (index) => {
        const updatedOptions = options.filter((_, i) => i !== index);
        setOptions(updatedOptions);
    };

    const handleSubmit = (e) => {
        e.preventDefault();

        if (!description.trim()) {
            setError("Poll description is required.");
            return;
        }

        if (options.some((option) => !option.option_text.trim())) {
            setError("All options must have text.");
            return;
        }

        const pollData = {
            description,
            options,
        };

        const requestOptions = {
            method: id ? "PUT" : "POST", 
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(pollData),
        };

        const url = id
            ? `${config.backendBaseUrl}/poll/${id}` 
            : `${config.backendBaseUrl}/poll`; 

        fetch(url, requestOptions)
            .then((response) => {
                if (!response.ok) {
                    throw new Error("Failed to save poll.");
                }
                return response.json();
            })
            .then(() => {
                navigate("/polls"); 
            })
            .catch((err) => {
                console.error(err);
                setError("An error occurred while saving the poll.");
            });
    };

    return (
        <div className="container mt-4">
            <h2>{id ? "Edit Poll" : "Add New Poll"}</h2>
            {error && <div className="alert alert-danger">{error}</div>}
            <form onSubmit={handleSubmit}>
                <div className="mb-3">
                    <label htmlFor="description" className="form-label">
                        Poll Description
                    </label>
                    <input
                        type="text"
                        id="description"
                        className="form-control"
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                        required
                    />
                </div>
                <div className="mb-3">
                    <label className="form-label">Options</label>
                    {options.map((option, index) => (
                        <div key={index} className="input-group mb-2">
                            <input
                                type="text"
                                className="form-control"
                                placeholder={`Option ${index + 1}`}
                                value={option.option_text}
                                onChange={(e) => handleOptionChange(index, e.target.value)}
                                required
                            />
                            <button
                                type="button"
                                className="btn btn-danger"
                                onClick={() => removeOption(index)}
                                disabled={options.length <= 1} // Prevent removing the last option
                            >
                                Remove
                            </button>
                        </div>
                    ))}
                    <button
                        type="button"
                        className="btn btn-primary"
                        onClick={addOption}
                    >
                        Add Option
                    </button>
                </div>
                <button type="submit" className="btn btn-success">
                    {id ? "Update Poll" : "Create Poll"}
                </button>
            </form>
        </div>
    );
};

export default EditPoll;