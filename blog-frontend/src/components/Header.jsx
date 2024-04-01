import React, { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";

const Header = () => {
  const userToken = localStorage.getItem("userToken") || sessionStorage.getItem("userToken");
  const navigate = useNavigate();
  const [userRole, setUserRole] = useState(null);
  const [userName, setUserName] = useState(null);

  const handleSignOut = () => {
    localStorage.removeItem("userToken");
    sessionStorage.removeItem("userToken");
    navigate("/signin");
  };

  useEffect(() => {
    try {
      if (userToken) {
        const decodedToken = jwtDecode(userToken);
        setUserRole(decodedToken.role);
        setUserName(decodedToken.name);
      }
    } catch (error) {
      console.error("Error decoding token:", error);
      navigate("/");
    }
  }, [userToken, navigate]);

  return (
    <div className="header">
      <div>
        {userToken ? (
          <Link className="logo" to="/">
            <span className="user-info">Welcome, {userName}</span>
          </Link>
        ) : (
          <Link className="logo" to="/">
            <span className="user-info">Blogs</span>
          </Link>
        )}
      </div>
      <div className="buttons">
        {userToken ? (
          <>
            {userRole === "blogger" && (
              <>
                <Link className="blogs" to="/my/blogs">
                  My Blogs
                </Link>
                <Link className="blogs" to="/add-blog">
                  Add Blog
                </Link>
              </>
            )}
            <button className="signout" onClick={handleSignOut}>
              Sign Out
            </button>
          </>
        ) : (
          <>
            <Link className="signup" to="/signup">
              Sign Up
            </Link>
            <Link className="signin" to="/signin">
              Sign In
            </Link>
          </>
        )}
      </div>
    </div>
  );
};

export default Header;
