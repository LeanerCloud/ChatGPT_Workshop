import React from 'react';
import BirthdayForm from '../components/BirthdayForm';

function HomePage() {
  return (
    <div>
      <h1>Home</h1>
      <BirthdayForm />
      {/* Render list of birthdays here... */}
    </div>
  );
}

export default HomePage;
