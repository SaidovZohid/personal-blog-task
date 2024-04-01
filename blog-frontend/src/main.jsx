import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.jsx';
import './index.css';

import SignUp from './pages/SignUp.jsx';
import SignIn from './pages/SignIn.jsx';
import AddBlog from './pages/AddBlog.jsx';
import SingleBlog from './pages/SingleBlog.jsx';
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import Header from './components/Header.jsx';
import { isAuthenticated } from './authUtils.js';
import MyBlogs from './pages/MyBlogs.jsx';
import UpdateBlogPage from './pages/UpdateBlogPage.jsx';

const ProtectedRoute = ({ element }) => {
  return isAuthenticated() ? element : <Navigate to="/signin" replace />;
};

const root = ReactDOM.createRoot(document.getElementById('root'));

root.render(
  <React.StrictMode>
    <BrowserRouter>
      <Header />
      <Routes>
      <Route
          path="/"
          element={<ProtectedRoute element={<App />} />}
        />
        <Route path="/signup" element={<SignUp />} />
        <Route path="/signin" element={<SignIn />} />
        <Route
          path="/add-blog"
          element={<ProtectedRoute element={<AddBlog />} />}
        />
        <Route path="/blogs/:id" element={<SingleBlog />} />
        <Route
          path="/update-blog/:id"
          element={<ProtectedRoute element={<UpdateBlogPage />} />}
        />
        <Route path="/my/blogs" element={<MyBlogs />} />
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
);
