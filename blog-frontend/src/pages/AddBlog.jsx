import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";

const apiBaseUrl = process.env.API_BASE_URL;

const AddBlog = () => {
  const navigate = useNavigate();
  const [header, setHeader] = useState("");
  const [body, setBody] = useState("");
  const [userRole, setUserRole] = useState(null);

  useEffect(() => {
    const token = localStorage.getItem("userToken") || sessionStorage.getItem("userToken");
    if (!token) {
      navigate("/");
      return;
    }
    try {
      const decodedToken = jwtDecode(token);
      setUserRole(decodedToken.role);
      if (decodedToken.role !== "blogger") {
        navigate("/");
      }
    } catch (error) {
      console.error("Invalid token:", error);
      navigate("/");
    }
  }, [navigate]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const token = localStorage.getItem("userToken") || sessionStorage.getItem("userToken");
      const response = await axios.post(
        `${apiBaseUrl}/posts`,
        { header, body },
        {
          headers: {
            Authorization: token,
          },
        }
      );
      navigate("/");
    } catch (error) {
      console.log(error)
      if (error.response.data.code === 401) {
        navigate("/");
      }
    }
  };

  return (
    <div className="add-blog-container">
      <h2>Add a New Blog Post</h2>
      {userRole !== "blogger" && <p>You are not authorized to create blog posts.</p>}
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="header" className="add-blog-label">
            Header:
          </label>
          <input className="add-blog-input" type="text" id="header" value={header} onChange={(e) => setHeader(e.target.value)} required />
        </div>
        <div className="form-group">
          <label htmlFor="body" className="add-blog-label">
            Body:
          </label>
          <textarea className="add-blog-input" id="body" value={body} onChange={(e) => setBody(e.target.value)} required />
        </div>
        <button type="submit" className="add-blog-header">
          Submit
        </button>
      </form>
    </div>
  );
};

export default AddBlog;
