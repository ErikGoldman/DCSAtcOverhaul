namespace ATCLib
{
  public abstract class MessagePayload
  {
    public string RawContent { get; }

    public MessagePayload(string rawContent)
    {
      RawContent = rawContent;
    }

    public static MessagePayload? Parse(string rawContent)
    {
      throw new NotImplementedException();
    }
  }

  public class MessageParseException : Exception
  {
    public MessageParseException(string message) : base(message) { }
  }

  public readonly struct MessageToken(bool isCallsign, string content)
  {
    public bool IsCallsign { get; } = isCallsign;
    public string Content { get; } = content;
  }

  public class Message
  {
    public string RawMessage { get; }

    public ICommunicator From { get; }
    public ICommunicator To { get; }
    public MessagePayload Payload { get; }
    public static List<Type> MessagePayloadTypes { get; } = [];

    public Message(string rawMessage, ICommunicator from, ICommunicator to, MessagePayload payload)
    {
      RawMessage = rawMessage;
      From = from;
      To = to;
      Payload = payload;
    }

    // this function tokenizes the message by extracting the callsigns and the message content separately
    public static List<MessageToken> SplitMessage(string rawMessage, ActiveCommunicatorList communicators)
    {
      if (rawMessage == "")
      {
        return [];
      }

      string[] inTokens = rawMessage.Split(' ');
      if (inTokens.Length == 0)
      {
        return [];
      }

      List<MessageToken> outTokens = [];
      int currentCallsignParsePositionStart = 0;
      int currentCallsignParsePositionEnd = 0;

      while (currentCallsignParsePositionStart < inTokens.Length)
      {
        var callsignSoFar = string.Join(" ", inTokens[currentCallsignParsePositionStart..(currentCallsignParsePositionEnd + 1)]);

        if (communicators.GetByName(callsignSoFar) != null)
        {
          outTokens.Add(new MessageToken(true, callsignSoFar));
          currentCallsignParsePositionStart = currentCallsignParsePositionEnd + 1;
          currentCallsignParsePositionEnd = currentCallsignParsePositionStart;
        }
        else if (communicators.PrefixExists(callsignSoFar))
        {
          currentCallsignParsePositionEnd++;
        }
        else
        {
          outTokens.Add(new MessageToken(false, inTokens[currentCallsignParsePositionStart]));
          currentCallsignParsePositionStart++;
          currentCallsignParsePositionEnd = currentCallsignParsePositionStart;
        }
      }

      return outTokens;
    }

    public static Message? Parse(string rawMessage, ATCState state)
    {
      /*
      string[] parts = rawMessage.Split(' ', 3);
      ICommunicator? maybeFrom = null;
      ICommunicator? maybeTo = null;
      string content;

      switch (parts.Length)
      {
        case 1:
          content = parts[0];
          break;
        case 2:
          maybeFrom = state.GetByName(parts[^1]);
          if (maybeFrom == null)
          {
            content = rawMessage;
          }
          else
          {
            content = parts[0];
          }
          break;
        case 3:
          maybeTo = state.GetByName(parts[0]);
          maybeFrom = state.GetByName(parts[1]);
          content = parts[2];
          break;
        default:
          throw new MessageParseException("Invalid message format");
      }

      if (maybeFrom == null && maybeTo == null)
      {
        throw new MessageParseException("Unable to determine sender or recipient");
      }

      MessagePayload? payload = null;

      foreach (var payloadType in MessagePayloadTypes)
      {
        var parseMethod = payloadType.GetMethod("Parse", System.Reflection.BindingFlags.Public | System.Reflection.BindingFlags.Static);
        if (parseMethod != null)
        {
          try
          {
            payload = parseMethod.Invoke(null, new object[] { content }) as MessagePayload;
            if (payload != null)
            {
              break;
            }
          }
          catch (Exception)
          {
            // Log the error if needed
          }
        }
      }

      if (payload == null)
      {
        return null;
      }

      return new Message(rawMessage, maybeFrom ?? state.GetByName("UNKNOWN"), maybeTo ?? state.GetByName("UNKNOWN"), payload);
      */
      throw new NotImplementedException();
    }
  }
}