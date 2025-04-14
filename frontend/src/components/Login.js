import { useState } from "react";
import { useNavigate, useOutletContext } from "react-router-dom";
import Input from "./form/Input";

const Login = () => {
    const [isRegister, setIsRegister] = useState(false);
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [username, setUsername] = useState("");

    const { setJwtToken, setUserId } = useOutletContext(); // Add setUserId to context
    const { setAlertClassName } = useOutletContext();
    const { setAlertMessage } = useOutletContext();
    const { toggleRefresh } = useOutletContext();

    const navigate = useNavigate();

    const handleSubmit = (event) => {
        event.preventDefault();

        let payload = {
            email: email,
            password: password,
        };

        if (isRegister) {
            payload.username = username;
        }

        const headers = new Headers();
        headers.append("Content-Type", "application/json");

        const requestOptions = {
            method: "POST",
            headers: headers,
            body: JSON.stringify(payload),
        };

        const url = isRegister
            ? `http://localhost:8080/v1/user`
            : `http://localhost:8080/v1/auth/token`;

        fetch(url, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                console.log("user", data);
                if (data.error) {
                    setAlertClassName("alert-danger");
                    setAlertMessage(data.message);
                } else {
                    if (!isRegister) {
                        setJwtToken(data.access_token);
                        setAlertClassName("d-none");
                        setAlertMessage("");
                        toggleRefresh(true);
                        navigate("/");
                    } else {
                        setAlertClassName("alert-success");
                        setAlertMessage("Registration successful! Please log in.");
                        setUserId(data.user_id); // Save userId in context
                        setIsRegister(false);
                    }
                }
            })
            .catch((error) => {
                setAlertClassName("alert-danger");
                setAlertMessage(error.message || "An error occurred.");
            });
    };

    return (
        <div className="col-md-6 offset-md-3">
            <h2>{isRegister ? "Register" : "Login"}</h2>
            <hr />

            <form onSubmit={handleSubmit}>
                {isRegister && (
                    <Input
                        title="Username"
                        type="text"
                        className="form-control"
                        name="username"
                        autoComplete="username-new"
                        onChange={(event) => setUsername(event.target.value)}
                    />
                )}

                <Input
                    title="Email Address"
                    type="email"
                    className="form-control"
                    name="email"
                    autoComplete="email-new"
                    onChange={(event) => setEmail(event.target.value)}
                />

                <Input
                    title="Password"
                    type="password"
                    className="form-control"
                    name="password"
                    autoComplete="password-new"
                    onChange={(event) => setPassword(event.target.value)}
                />

                <hr />

                <input
                    type="submit"
                    className="btn btn-primary"
                    value={isRegister ? "Register" : "Login"}
                />
            </form>

            <hr />
            <p className="text-center">
                {isRegister ? (
                    <>
                        Already registered?{" "}
                        <button
                            type="button"
                            className="btn btn-link"
                            onClick={() => setIsRegister(false)}
                        >
                            Login here
                        </button>
                    </>
                ) : (
                    <>
                        Don't have an account?{" "}
                        <button
                            type="button"
                            className="btn btn-link"
                            onClick={() => setIsRegister(true)}
                        >
                            Register here
                        </button>
                    </>
                )}
            </p>
        </div>
    );
};

export default Login;