import React, { useState } from 'react';
import axios from 'axios';

function GroupForm() {
  const [groupName, setGroupName] = useState("");
  const [operation, setOperation] = useState("join"); // either "join" or "leave"

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (operation === "join") {
      await axios.post('/group/join', { name: groupName });
    } else if (operation === "leave") {
      await axios.post('/group/leave', { name: groupName });
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <label>
        Group Name:
        <input 
          type="text" 
          value={groupName} 
          onChange={(e) => setGroupName(e.target.value)} 
          required
        />
      </label>

      <div>
        <label>
          <input 
            type="radio" 
            value="join" 
            checked={operation === "join"}
            onChange={() => setOperation("join")}
          />
          Join
        </label>

        <label>
          <input 
            type="radio" 
            value="leave" 
            checked={operation === "leave"}
            onChange={() => setOperation("leave")}
          />
          Leave
        </label>
      </div>

      <button type="submit">Submit</button>
    </form>
  );
}

export default GroupForm;
