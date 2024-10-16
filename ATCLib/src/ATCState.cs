namespace ATCLib
{
  public class ActiveCommunicatorList
  {
    Dictionary<string, ICommunicator> Communicators { get; }

    public ActiveCommunicatorList()
    {
      Communicators = [];
    }

    public void AddCommunicator(ICommunicator communicator)
    {
      Communicators.Add(communicator.GetCallsign(), communicator);
    }

    public ICommunicator? GetByName(string name)
    {
      return Communicators.GetValueOrDefault(name);
    }

    public bool PrefixExists(string prefix)
    {
      return Communicators.Any(c => c.Key.StartsWith(prefix));
    }
  }

  public class ATCState
  {
    ActiveCommunicatorList ActiveCommunicators { get; }

    public ATCState()
    {
      ActiveCommunicators = new ActiveCommunicatorList();
    }
  }
}