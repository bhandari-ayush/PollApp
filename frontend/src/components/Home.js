import { Link } from "react-router-dom";
import Ticket from "./../images/poll.jpg";

const Home = () => {
  return (
    <>
      <div className="text-center">
        <h2> Let's Poll!!</h2>
        <hr />
        <Link to="/polls">
          <img src={Ticket} alt="movie tickets"></img>
        </Link>
      </div>
    </>
  );
};

export default Home;
