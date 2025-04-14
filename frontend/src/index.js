import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import App from './App';
import EditMovie from './components/EditMovie';
import ErrorPage from './components/ErrorPage';
import Home from './components/Home';
import Login from './components/Login';
import Polls from './components/Polls';
import Poll from './components/Poll';
import VoteResult from './components/VoteResult';
import EditPoll from './components/EditPoll';

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    errorElement: <ErrorPage />,
    children: [
      {index: true, element: <Home /> },
      {
        path: "/admin/movie/0",
        element: <EditMovie />,
      },
      {
        path: "/poll/0",
        element: <EditPoll />,
      },
      {
        path: "/polls",
        element: <Polls />,
      },
      {
        path: "/poll/:id",
        element: <Poll />,
      },
      {
        path: "/vote/:id",
        element: <VoteResult />,
      },
      {
        path: "/login",
        element: <Login />,
      },
    ]
  }
])

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
