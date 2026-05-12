const TourCard = ({ name = 'Tour Name', price = '0đ' }) => {
  return (
    <article style={{ border: '1px solid #ddd', borderRadius: 12, padding: 16 }}>
      <h3>{name}</h3>
      <p>{price}</p>
    </article>
  );
};

export default TourCard;
