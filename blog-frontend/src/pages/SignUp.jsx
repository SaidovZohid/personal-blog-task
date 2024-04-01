import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate } from 'react-router-dom';
import { jwtDecode } from "jwt-decode";

const apiBaseUrl = process.env.API_BASE_URL;

const SignUp = () => {
  const navigate = useNavigate();
  useEffect(() => {
    const userToken = localStorage.getItem('userToken') || sessionStorage.getItem('userToken');
    if (userToken) {
      navigate('/');
    }
  }, []);

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState("");
  const [code, setCode] = useState("");
  const [isSignedUp, setIsSignedUp] = useState(false);
  const [error, setError] = useState("");

  const handleSignUp = async () => {
    const payload = {
      email,
      password,
      role,
    };
    try {
      const response = await axios.post(`${apiBaseUrl}/auth/signup`, payload);
      console.log(response.data);
      setIsSignedUp(true);
    } catch (error) {
      console.log(error)
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

  const handleVerification = async () => {
    const payload = {
      email,
      code,
    };
    try {
      const response = await axios.post(`${apiBaseUrl}/auth/verify`, payload);
      const token = response.data.access_token;
      localStorage.setItem("userToken", token);
      if (token) {
        const decodedToken = jwtDecode(token);
        if (decodedToken.role == "blogger") {
          navigate("/my/blogs");
        } else {
          navigate('/');
        }
      }
    } catch (error) {
      if (error.response) {
        if (error.response.status === 400) {
          setError("Invalid email or password.");
        } else if (error.response.status === 500) {
          setError("Internal server error. Please try again later.");
        } else if (error.response.status == 404) {
          setError("User not found");
        } else {
          setError("An unexpected error occurred.");
        }
      } else {
        setError("Network error. Please check your internet connection.");
      }
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!isSignedUp) {
      handleSignUp();
    } else {
      handleVerification();
    }
  };

  const handleRoleChange = (e) => {
    setRole(e.target.value);
  };


  return (
    <div>
      {error && <div className="error">{error}</div>}
      <form className="signup-form" onSubmit={handleSubmit}>
        <label>Email:</label>
        <input type="email" placeholder="youremail@gmail.com" onChange={(e) => setEmail(e.target.value)} value={email} disabled={isSignedUp} required/>

        <label>Password:</label>
        <input type="password" placeholder="********" onChange={(e) => setPassword(e.target.value)} value={password} disabled={isSignedUp} required/>

        {isSignedUp && (
          <>
            <label>Verification Code:</label>
            <input type="text" placeholder="Enter the 6 digit code" onChange={(e) => setCode(e.target.value)} value={code} required/>
          </>
        )}

        <label>Select Your Role:</label>
        <div className="radio">
          <input type="radio" name="role" id="blogger" value="blogger" onChange={handleRoleChange} disabled={isSignedUp} />
          <label htmlFor="blogger">Blogger</label>

          <input type="radio" name="role" id="reader" value="reader" onChange={handleRoleChange} disabled={isSignedUp} />
          <label htmlFor="reader">Reader</label>
        </div>
        <button type="submit" className="submit-button-signin-signup">{isSignedUp ? "Verify" : "Sign Up"}</button>
      </form>
    </div>
  );
};

export default SignUp;
