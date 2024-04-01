import React, { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import axios from "axios";
import { jwtDecode } from "jwt-decode";

const apiUrl = process.env.API_BASE_URL;

const UpdateBlogPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [header, setHeader] = useState("");
  const [body, setBody] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchPost = async () => {
      const token = localStorage.getItem("userToken") || sessionStorage.getItem("userToken");
      if (!token) {
        navigate("/");
        return;
      }
      try {
        const decodedToken = jwtDecode(token);
        if (decodedToken.role !== "blogger") {
          navigate("/");
        }


        const response = await axios.get(`${apiUrl}/posts/${id}`);
        setHeader(response.data.header);
        setBody(response.data.body);

        if (response.data.user_id !== decodedToken.user_id) {
          navigate("/");
        }
      } catch (error) {
        console.error("Invalid token:", error);
        navigate("/");
      }
    };
    fetchPost();
  }, [navigate, id]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const token = localStorage.getItem("userToken") || sessionStorage.getItem("userToken");
      if (!token) {
        throw new Error("Authorization token not found.");
      }

      const config = {
        headers: {
          Authorization: token,
        },
      };

      const requestData = {
        header,
        body,
      };

      const response = await axios.put(`${apiUrl}/posts/${id}`, requestData, config);
      console.log("Post updated:", response.data);
      navigate(`/blogs/${id}`);
    } catch (error) {
      console.error("Error updating post:", error);
      setError("Failed to update post. Please try again.");
    }
  };

  return (
    <div className="update-blog-container">
      {error && <p className="error">{error}</p>}
      <h2>Update Your Blog</h2>
      <form onSubmit={handleSubmit}>
        <input type="text" placeholder="Enter new header" value={header} onChange={(e) => setHeader(e.target.value)} className="update-blog-input" />
        <textarea placeholder="Enter new body" value={body} onChange={(e) => setBody(e.target.value)} className="update-blog-textarea"></textarea>
        <button type="submit" className="update-blog-button">
          Update Post
        </button>
      </form>
    </div>
  );
};

export default UpdateBlogPage;
