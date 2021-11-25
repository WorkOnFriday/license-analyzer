package work;

import pri.PriInt;
import pri.PriStat;
import dp2.Dp2PriStm;
import dp3.Dp3PriStm;

public class Main {

    public static void main(String[] args) {
        PriInt pi = new PriInt();
        pi.priInt();
        PriStat ps = new PriStat();
        ps.priStat();
        Dp2PriStm dp2 = new Dp2PriStm();
        dp2.priStm();
        Dp3PriStm dp3 = new Dp3PriStm();
        dp3.priStm();
    }
}
