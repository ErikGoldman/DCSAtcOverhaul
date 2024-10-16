namespace ATCLib
{
  public class Aircraft : ICommunicator
  {
    public string Type { get; }
    public string Callsign { get; }

    public Aircraft(string callsign, string type)
    {
      Callsign = callsign;
      Type = type;
    }

    public override string ToString()
    {
      return $"{((ICommunicator)this).GetCallsign()} ({Type})";
    }

    void ICommunicator.SendMessage(Message message)
    {
      throw new NotImplementedException();
    }

    void ICommunicator.ReceiveMessage(Message message)
    {
      throw new NotImplementedException();
    }

    string ICommunicator.GetCallsign()
    {
      return Callsign;
    }
  }
}