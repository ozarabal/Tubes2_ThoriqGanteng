import React, {useState, useEffect} from 'react';
import { List, ListItem, Card } from '@material-tailwind/react';



const ListDefault = ({ data }) => {
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (data && data.result && Array.isArray(data.result) && data.result.length > 0) {
      setIsLoading(false);
    }
  }, [data]);

  if (isLoading) {
    return <div>Loading...</div>; // Tampilkan pesan loading jika data masih belum siap
  }

  if (!data || !data.result || !Array.isArray(data.result) || data.result.length === 0) {
    return null; // Kembalikan null jika tidak ada data atau data tidak sesuai
  }

  const resultList = data.result;

  return (
    <div className="flex flex-col items-center justify-center">
      <Card className="w-96">
        <List>
          {resultList.map((item, index) => (
            <ListItem key={index}>
              <a href={item} target="_blank" rel="noopener noreferrer">
                {item}
              </a>
            </ListItem>
          ))}
        </List>
      </Card>
    </div>
  );
}

export default ListDefault;
