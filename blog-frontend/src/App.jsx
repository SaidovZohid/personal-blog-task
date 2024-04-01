import React, { useState, useEffect } from "react";
import axios from "axios";
import { Link } from "react-router-dom";

const apiBaseUrl = process.env.API_BASE_URL;

const App = () => {
  const [posts, setPosts] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [postCount, setPostCount] = useState(0);

  const fetchPosts = async (page) => {
    try {
      const response = await axios.get(`${apiBaseUrl}/posts?page=${page}&limit=10&sort=desc`);
      setPosts(response.data.posts);
      setTotalPages(Math.ceil(response.data.count / 10));
      setPostCount(response.data.count);
    } catch (error) {
      console.error("Error fetching posts:", error);
    }
  };

  useEffect(() => {
    fetchPosts(currentPage);
  }, [currentPage]);

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
      <h1>Blogs</h1>
      <div className="all-posts-container">
        {postCount === 0 ? <p>No Blogs :(</p> : null}
        {posts.map((post) => (
          <Link key={post.id} to={`/blogs/${post.id}`} className="post-link">
            <div className="post-card">
              <h3 className="post-header">{post.header.substring(0, 100)}
                {post.header.length > 300 ? "..." : ""}</h3>
              <p className="post-body">
                {post.body.substring(0, 100)}
                {post.body.length > 300 ? "..." : ""}
              </p>
              <p className="post-author">Author: {post.user_info.name}</p>
              <p className="post-date">{formatDate(post.created_at)}</p>
            </div>
          </Link>
        ))}
      </div>
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
  );
};

export default App;
