import { Link } from "react-router-dom";
import PollImg from "./../images/poll.jpg";

const Home = () => {
  return (
    <>
      <div className="text-center">
        <h2> Let's Poll!!</h2>
        <hr />
        <Link to="/polls">
          <img src={PollImg} alt="Let's Poll"></img>
        </Link>
      </div>
    </>
  );
};

export default Home;
