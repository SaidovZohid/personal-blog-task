import axios from "axios";
import { jwtDecode } from "jwt-decode";
import React, { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { formatDistanceToNow } from "date-fns";

const apiUrl = process.env.API_BASE_URL;

const SingleBloog = () => {
  const { id } = useParams();
  const [post, setPost] = useState(null);
  const [loading, setLoading] = useState(true);
  const [addComment, setAddComment] = useState(false);
  const [commentText, setCommentText] = useState("");
  const [shouldSignIn, setShoudlSignIn] = useState("");
  const [loggedInUserInfo, setLoggedInUserInfo] = useState("");
  const [replyText, setReplyText] = useState("");
  const [replyToComment, setReplyToComment] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchPost = async (postId) => {
      try {
        const response = await axios.get(`${apiUrl}/posts/${postId}`);
        setPost(response.data);
      } catch (error) {
        console.error("Error fetching post:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchPost(id);
  }, [id]);

  useEffect(() => {
    try {
      const token = localStorage.getItem("userToken") || sessionStorage.getItem("userToken");
      if (token) {
        const decodedToken = jwtDecode(token);
        setLoggedInUserInfo(decodedToken);
      }
    } catch (error) {
      console.log(error);
      if (error.response.data.code === 401) {
        navigate("/");
        return;
      }
    }
  }, [navigate]);

  const handleAddComment = () => {
    try {
      const token = localStorage.getItem("userToken") || sessionStorage.getItem("userToken");
      if (!token) {
        throw new Error("To add comment. You should be signed in!");
      }
      const decodedToken = jwtDecode(token);
    } catch (error) {
      console.log(error);
      setShoudlSignIn("To add comment. You should be signed in!");
      return;
    }
    console.log("Here is it");
    setAddComment(true);
  };

  const handleSubmit = () => {
    try {
      if (commentText.trim() === "") {
        return;
      }
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
        content: commentText,
        post_id: parseInt(id),
      };
      const response = axios.post(`${apiUrl}/comments`, requestData, config);
      console.log("Comment created:", response.data);
    } catch (error) {
      console.error("Invalid token:", error);
      navigate("/");
      return;
    }
  };

  const handleDeleteComment = async (commentId) => {
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
      await axios.delete(`${apiUrl}/comments/${commentId}`, config);

      const updatedComments = post.all_comments.comments.filter((comment) => comment.id !== commentId);
      setPost((prevPost) => ({
        ...prevPost,
        all_comments: {
          ...prevPost.all_comments,
          comments: updatedComments,
        },
      }));
      post.all_comments.count = post.all_comments.count - 1;
    } catch (error) {
      console.error("Invalid token:", error);
      navigate("/");
      return;
    }
  };

  const handleDeletePost = async (postId) => {
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
      await axios.delete(`${apiUrl}/posts/${postId}`, config);
      navigate("/");
    } catch (error) {
      console.error("Invalid token:", error);
      navigate("/");
      return;
    }
  };

  const handleReplyToComment = async (commentId) => {
    if (!loggedInUserInfo) {
      setShoudlSignIn("Please sign in or sign up to add a reply to this comment.");
      return;
    }
    // Please sign in or sign up to add a comment to this blog.
    setReplyToComment(commentId);
  };

  const formatDate = (timestamp) => {
    const date = new Date(timestamp);
    const options = { month: "short", day: "numeric", year: "numeric" };
    return date.toLocaleDateString("en-US", options);
  };

  const getTheDurationTime = (timestamp) => {
    const date = new Date(timestamp);
    const duration = formatDistanceToNow(date, { addSuffix: true });
    return duration;
  };

  const handleCancel = () => {
    setAddComment(false);
  };

  const handleCancelReply = () => {
    setReplyToComment(false);
  };

  const handleReplySubmit = (event, comment_id) => {
    event.preventDefault();
    if (replyText.trim() === "") {
      return;
    }
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
      content: replyText,
      comment_id: comment_id,
      post_id: parseInt(id),
    };

    axios
      .post(`${apiUrl}/replies`, requestData, config)
      .then((response) => {
        handleCancelReply();
        console.log("Reply created:", response.data);
        const newReply = response.data;
        setPost((prevPost) => {
          const updatedComments = prevPost.all_comments.comments.map((comment) => {
            if (comment.id === comment_id) {
              const all_replies = comment.all_replies ? { ...comment.all_replies } : { replies: [], count: 0 };
              return {
                ...comment,
                all_replies: {
                  ...all_replies,
                  replies: [...all_replies.replies, newReply],
                  count: all_replies.count + 1,
                },
              };
            }
            return comment;
          });
          return {
            ...prevPost,
            all_comments: {
              ...prevPost.all_comments,
              comments: updatedComments,
            },
          };
        });
      })
      .catch((error) => {
        setShoudlSignIn(error);
      });
    setReplyText("");
  };

  const handleDeleteReply = async (replyId) => {
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
      await axios.delete(`${apiUrl}/replies/${replyId}`, config);

      if (post.all_comments && post.all_comments.comments) {
        const updatedComments = post.all_comments.comments.map((comment) => {
          if (comment.all_replies && comment.all_replies.count > 0) {
            if (comment.all_replies.replies.some((reply) => reply.id === replyId)) {
              const updatedReplies = comment.all_replies.replies.filter((reply) => reply.id !== replyId);
              return {
                ...comment,
                all_replies: {
                  ...comment.all_replies,
                  replies: updatedReplies,
                  count: comment.all_replies.count - 1,
                },
              };
            }
          }
          return comment;
        });
        setPost((prevPost) => ({
          ...prevPost,
          all_comments: {
            ...prevPost.all_comments,
            comments: updatedComments,
          },
        }));
      }
    } catch (error) {
      console.error("Invalid token:", error);
      navigate("/");
      return;
    }
  };

  return (
    <div className="single-blog-container">
      {loading ? (
        <div className="single-loading">
          <div className="single-loader"></div>
          <p>Loading...</p>
        </div>
      ) : (
        post && (
          <div>
            {loggedInUserInfo && loggedInUserInfo.user_id === post.user_id && (
              <div className="single-post-buttons">
                <button className="single-delete-button" onClick={() => handleDeletePost(post.id)}>
                  Delete
                </button>
                <Link className="single-update-link" to={`/update-blog/${post.id}`}>
                  Update Blog
                </Link>
              </div>
            )}
            <h1 className="single-post-header">{post.header}</h1>
            <p className="single-author-name">
              Author: {post.user_info.name} - {formatDate(post.created_at)}
            </p>
            <p className="single-post-body">{post.body}</p>
            <div className="single-comments-section">
              <button className="add-comment-button" onClick={handleAddComment}>
                Add Comment
              </button>
              {shouldSignIn !== "" ? <p className="single-should-signin-up">{shouldSignIn}</p> : null}
              <h2 className="comments-count">{post.all_comments.count} Comments:</h2>
              {addComment ? (
                <form>
                  <input type="text" className="single-text-area" placeholder="Add a comment..." required value={commentText} onChange={(e) => setCommentText(e.target.value)} />
                  <div className="single-two-buttons-cancel-submit">
                    <button type="submit" onClick={handleSubmit} className="add-comment-button">
                      Comment
                    </button>
                    <button type="submit" onClick={handleCancel} className="add-comment-button">
                      Cancel
                    </button>
                  </div>
                </form>
              ) : null}
              {post.all_comments.count === 0 ? (
                <p className="single-no-comments">No Comments</p>
              ) : (
                <>
                  <ul>
                    {post.all_comments.comments.map((comment) => (
                      <li key={comment.id} className="single-comment">
                        <p className="comment-date-format">
                          <span className="single-post-comment-username">@{comment.user_info.name}</span> {getTheDurationTime(comment.created_at)}
                        </p>
                        <p className="single-comment-content">{comment.content}</p>
                        <div className="comment-buttons-and-replies">
                          <div className="comment-buttons">
                            <>
                              {loggedInUserInfo && loggedInUserInfo.user_id === comment.user_id && (
                                <button className="single-delete-button" onClick={() => handleDeleteComment(comment.id)}>
                                  Delete
                                </button>
                              )}
                              <button className="single-reply-button" onClick={() => handleReplyToComment(comment.id)}>
                                Reply
                              </button>
                            </>
                          </div>
                          <p className="replies-count">Replies: {comment.all_replies.count}</p>
                        </div>
                        {replyToComment == comment.id ? (
                          <form
                            onSubmit={(e) => {
                              e.preventDefault();
                              handleReplySubmit(e, comment.id);
                            }}
                            className="single-add-reply-to-comment">
                            <input type="text" className="single-text-area" placeholder="Add a reply to comment..." required value={replyText} onChange={(e) => setReplyText(e.target.value)} />
                            <div className="single-add-reply-buttons">
                              <button disabled={replyText === ""} type="submit" className="single-add-reply-button">
                                Reply
                              </button>
                              <button type="cancel" onClick={handleCancelReply} className="single-add-reply-cancel-button">
                                Cancel
                              </button>
                            </div>
                          </form>
                        ) : null}
                        {comment.all_replies.count === 0 ? null : (
                          <ol className="single-comment-replies">
                            {comment.all_replies.replies.map((reply) => (
                              <li key={reply.id} className="single-comment-reply">
                                <p className="comment-date-format">
                                  <span className="single-post-comment-username">@{reply.user_info.name}</span> {getTheDurationTime(reply.created_at)}
                                </p>
                                <p className="single-comment-content">{reply.content}</p>
                                <div className="comment-buttons">
                                  <>
                                    {loggedInUserInfo && loggedInUserInfo.user_id === reply.user_id && (
                                      <button className="single-delete-button" onClick={() => handleDeleteReply(reply.id)}>
                                        Delete
                                      </button>
                                    )}
                                  </>
                                </div>
                              </li>
                            ))}
                          </ol>
                        )}
                      </li>
                    ))}
                  </ul>
                </>
              )}
            </div>
          </div>
        )
      )}
    </div>
  );
};

export default SingleBloog;
