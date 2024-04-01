import React, { useState, useEffect } from "react";
import axios from "axios";
import { Link, useNavigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";

const apiUrl = process.env.API_BASE_URL;

const MyBlogs = () => {
  const [posts, setPosts] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [userId, setUserId] = useState(-1);
  const [blogsCount, setBlogsCount] = useState();
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("userToken") || sessionStorage.getItem("userToken");
    if (!token) {
      navigate("/");
      return;
    }
    try {
      const decodedToken = jwtDecode(token);
      setUserId(decodedToken.user_id);
      if (decodedToken.role !== "blogger") {
        navigate("/");
        return;
      }
      fetchPosts(currentPage, decodedToken.user_id);
    } catch (error) {
      console.error("Invalid token:", error);
      navigate("/");
      return;
    }
  }, [currentPage, navigate]);

  const fetchPosts = async (page, userId) => {
    try {
      const response = await axios.get(`${apiUrl}/posts?page=${page}&limit=10&sort=desc&user_id=${userId}`);
      setPosts(response.data.posts);
      setBlogsCount(response.data.count);
      setTotalPages(Math.ceil(response.data.count / 10));
    } catch (error) {
      console.error("Error fetching posts:", error);
    }
  };

  const handleNextPage = () => {
    if (currentPage < totalPages) {
      setCurrentPage(currentPage + 1);
    }
  };

  const handlePrevPage = () => {
    if (currentPage > 1) {
      setCurrentPage(currentPage - 1);
    }
  };

  const formatDate = (timestamp) => {
    const date = new Date(timestamp);
    const options = { month: 'short', day: 'numeric', year: 'numeric' };
    return date.toLocaleDateString('en-US', options);
  };
  

  return (
    <div className="all-posts">
      <h1>Your Blogs</h1>
      <div className="all-posts-container">
        {blogsCount === 0 ? <p>No Blogs</p> : null}
        {posts.map((post) => (
          <Link key={post.id} to={`/blogs/${post.id}`} className="post-link">
            <div className="post-card">
              <h3 className="post-header">{post.header}</h3>
              <p className="post-body">
                {post.body.substring(0, 100)}
                {post.body.length > 300 ? "..." : ""}
              </p>
              <p className="post-author">Author: {post.user_info.name}</p>
              <p className="post-date">{formatDate(post.created_at)}</p>
            </div>
          </Link>
        ))}
        {totalPages > 1 && (
          <div className="pagination">
            <button onClick={handlePrevPage} disabled={currentPage === 1}>
              Previous
            </button>
            <button onClick={handleNextPage} disabled={currentPage === totalPages}>
              Next
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default MyBlogs;
