namespace ATCLib
{
  public interface ICommunicator
  {
    void SendMessage(Message message);
    void ReceiveMessage(Message message);
    string GetCallsign();
  }

  public class Comms
  {
    public static void ReceiveRawMessage(string message, ATCState state)
    {
      Log.WriteLine("RAW << " + message);

      try
      {

      }
      catch (Exception e)
      {
        Log.WriteLine("Error: " + e.Message);
      }
    }
  }
}