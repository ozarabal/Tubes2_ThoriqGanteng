import React from 'react';
import { List, ListItem, Card } from '@material-tailwind/react';

const data = ['Item 1', 'Item 2', 'Item 3', 'Item 4'];

export default function ListFromData() {
  return (
    <div className="flex flex-col items-center justify-center">
      <Card className="w-96">
        <List>
          {data.map((item, index) => (
            <ListItem key={index}>{item}</ListItem>
          ))}
        </List>
      </Card>
    </div>
  );
}