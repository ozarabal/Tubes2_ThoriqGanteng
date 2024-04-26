import React, { useState } from 'react';
import { ToggleButton, ToggleButtonGroup } from '@mui/material';

export default function ColorToggleButton() {
  const [alignment, setAlignment] = useState('web');

  const handleChange = async (event, newAlignment) => {
    setAlignment(newAlignment);

    // Kirim data ke backend
    try {
      const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ tipe: newAlignment }), // Mengirim data ke backend
      };
      const response = await fetch('http://localhost:8080/submitAlignment', requestOptions);
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      const data = await response.json();
      console.log('Response from server:', data);
    } catch (error) {
      console.error('There was an error!', error);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center">
    <ToggleButtonGroup
      color="primary"
      value={alignment}
      exclusive
      onChange={handleChange}
      aria-label="Platform"
    >
      <ToggleButton value="BFS">BFS</ToggleButton>
      <ToggleButton value="IDS">IDS</ToggleButton>
    </ToggleButtonGroup>
    </div>
  );
}
