import React, { useState } from 'react';
import { ToggleButton, ToggleButtonGroup } from '@mui/material';

export default function ColorToggleButton2() {
  const [method, setMethod] = useState('web');

  const handleChange = async (event, newMethod) => {
    setMethod(newMethod);

    // Kirim data ke backend
    try {
      const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ tipe: newMethod }), // Mengirim data ke backend
      };
      const response = await fetch('http://localhost:8080/submitmethod', requestOptions);
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
      value={method}
      exclusive
      onChange={handleChange}
      aria-label="Platform"
    >
      <ToggleButton value="First">First</ToggleButton>
      <ToggleButton value="All">All</ToggleButton>
    </ToggleButtonGroup>
    </div>
  );
}