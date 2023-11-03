import React from 'react';
import logo from './logo.svg';
import HomePage from './views/HomePage';
import GroupPage from './views/GroupPage';
import './App.css';

function App() {
  // Assume a basic routing (for simplicity, not using React Router)
  let view;
  const path = window.location.pathname;
  if (path === '/group') {
    view = <GroupPage />;
  } else {
    view = <HomePage />;
  }

  return (
    <div className="App">
      {view}
    </div>
  );
}

export default App;