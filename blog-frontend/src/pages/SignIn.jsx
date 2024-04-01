import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";

const apiBaseUrl = process.env.API_BASE_URL;

const SignIn = () => {
  const navigate = useNavigate();
  useEffect(() => {
    const userToken = localStorage.getItem("userToken") || sessionStorage.getItem("userToken");
    if (userToken) {
      navigate("/");
    }
  }, []);

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [rememberMe, setRememberMe] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    const payload = {
      email,
      password,
      remember_me: rememberMe,
    };
    try {
      const response = await axios.post(`${apiBaseUrl}/auth/login`, payload);
      const token = response.data.access_token;
      if (rememberMe) {
        localStorage.setItem("userToken", token);
      } else {
        sessionStorage.setItem("userToken", token);
      }
      try {
        const decodedToken = jwtDecode(token);
        if (decodedToken.role == "blogger") {
          navigate("/my/blogs");
        } else {
          navigate("/");
        }
      } catch (error) {
        navigate("/signin")
      }
    } catch (error) {
      console.error("There was an error signing in:", error.response);
      if (error.response) {
        if (error.response.status === 400) {
          setError(error.response.data.message);
        } else if (error.response.status === 500) {
          setError("Internal server error. Please try again later.");
        } else {
          setError("An unexpected error occurred.");
        }
      } else {
        setError("Network error. Please check your internet connection.");
      }
    }
  };

  return (
    <div>
      {error && <div className="error">{error}</div>}
      <form className="signup-form" onSubmit={handleSubmit}>
        <label>Email:</label>
        <input type="email" placeholder="youremail@gmail.com" onChange={(e) => setEmail(e.target.value)} value={email} required />

        <label>Password:</label>
        <input type="password" placeholder="********" onChange={(e) => setPassword(e.target.value)} value={password} required />

        <label>
          <input type="checkbox" onChange={(e) => setRememberMe(e.target.checked)} checked={rememberMe} />
          Remember Me
        </label>

        <button type="submit" className="submit-button-signin-signup">
          Sign In
        </button>
      </form>
    </div>
  );
};

export default SignIn;
